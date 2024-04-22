import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChartDoughnutCitiesComponent } from './chart-doughnut-cities.component';

describe('ChartDoughnutCityesComponent', () => {
  let component: ChartDoughnutCitiesComponent;
  let fixture: ComponentFixture<ChartDoughnutCitiesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChartDoughnutCitiesComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ChartDoughnutCitiesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
